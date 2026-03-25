use std::fmt;

use tracing::{Event, Subscriber};
use tracing_subscriber::{
    fmt::{FmtContext, FormatEvent, FormatFields},
    registry::LookupSpan,
};

use time::{OffsetDateTime, UtcOffset};

pub struct TimeNodeFormatter;

impl<S, N> FormatEvent<S, N> for TimeNodeFormatter
where
    S: Subscriber + for<'a> LookupSpan<'a>,
    N: for<'a> FormatFields<'a> + 'static,
{
    fn format_event(
        &self,
        ctx: &FmtContext<'_, S, N>,
        mut writer: tracing_subscriber::fmt::format::Writer<'_>,
        event: &Event<'_>,
    ) -> fmt::Result {
        let jst = UtcOffset::from_hms(9, 0, 0).unwrap();

        let now = OffsetDateTime::now_utc().to_offset(jst);

        let fmt = time::format_description::parse("[year]/[month]/[day] [hour]:[minute]:[second]")
            .unwrap();

        write!(writer, "[{}] ", now.format(&fmt).unwrap())?;

        ctx.field_format().format_fields(writer.by_ref(), event)?;

        writeln!(writer)
    }
}
